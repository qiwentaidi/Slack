const MenuList = [
    {
      name: "aside.home",
      path: "/",
      icon: "/home.svg",
    },
    {
      name: "aside.penetration",
      path: "/Penetration", // 什么路径都行，但是别重复
      icon: "Smoking",
      children: [
        {
          name: "aside.webscan",
          path: "/Permeation/Webscan",
        },
        {
          name: "aside.portscan",
          path: "/Permeation/Portscan",
        },
        {
          name: "aside.dirscan",
          path: "/Permeation/Dirsearch",
        },
        {
          name: "aside.jsfinder",
          path: "/Permeation/JSFinder",
        },
        {
          name: "aside.exp",
          path: "/Permeation/Exploitation",
        },
      ]
    },
    {
      name: "aside.asset_collection",
      path: "/Asset",
      icon: "OfficeBuilding",
      children: [
        {
          name: "aside.asset_from_company",
          path: "/Asset/Company",
        },
        {
          name: "aside.subdomain_brute_force",
          path: "/Asset/Subdomain",
        },
        {
          name: "aside.dirscan",
          path: "/Asset/search_domain_info",
        },
      ]
    },
    {
      name: "aside.space_engine",
      path: "/SpaceEngine",
      icon: "Monitor",
      children: [
        {
          name: "aside.fofa",
          path: "/SpaceEngine/FOFA",
        },
        {
          name: "aside.hunter",
          path: "/SpaceEngine/Hunter",
        },
        {
          name: "aside.agent_pool",
          path: "/Asset/AgentPool",
        },
      ]
    },
    {
      name: "aside.tools",
      path: "/Tools",
      icon: "Tools",
      children: [
        {
          name: "aside.en_and_de",
          path: "/Tools/Codec",
        },
        {
          name: "aside.systeminfo",
          path: "/Tools/System",
        },
        {
          name: "aside.data_handing",
          path: "/Tools/DataHanding",
        },
        {
          name: "aside.memorandum",
          path: "/Tools/Memo",
        },
        {
          name: "aside.associate_dictionary_generator",
          path: "/Tools/Thinkdict",
        },
        {
          name: "aside.aksk",
          path: "/Tools/AKSK",
        },
      ]
    },
  ]

  export default MenuList