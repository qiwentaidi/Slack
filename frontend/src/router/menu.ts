const MenuList = [
    {
      name: "aside.home",
      path: "/",
      icon: "/home.svg",
    },
    {
      name: "aside.penetration",
      path: "/Penetration", // 父级路径怎么写都行，但是别重复
      icon: "Smoking",
      children: [
        {
          name: "aside.webscan",
          path: "/Permeation/Webscan",
          icon: "/webscan.png"
        },
        {
          name: "aside.portscan",
          path: "/Permeation/Portscan",
          icon: "/portscan.png"
        },
        {
          name: "aside.crack",
          path: "/Permeation/Crack",
          icon: "/target.png"
        },
        {
          name: "aside.dirscan",
          path: "/Permeation/Dirsearch",
          icon: "/dirscan.png"
        },
        {
          name: "aside.jsfinder",
          path: "/Permeation/JSFinder",
          icon: "/extract.png"
        },
        {
          name: "aside.exp",
          path: "/Permeation/Exploitation",
          icon: "/attack.png"
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
          icon: "/company.png",
        },
        {
          name: "aside.subdomain_brute_force",
          path: "/Asset/Subdomain",
          icon: "/bomb.png"
        },
        {
          name: "aside.search_domain_info",
          path: "/Asset/Ipdomain",
          icon: "/domain.png"
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
          icon: "/fofa.ico"
        },
        {
          name: "aside.hunter",
          path: "/SpaceEngine/Hunter",
          icon: "/hunter.ico"
        },
        {
          name: "aside.360quake",
          path: "/SpaceEngine/Quake",
          icon: "/navigation/360.ico"
        },
        {
          name: "aside.polymerization",
          path: "/SpaceEngine/Polymerization",
          icon: "/navigation/polymerization.png"
        },
        {
          name: "aside.agent_pool",
          path: "/SpaceEngine/AgentPool",
          icon: "/navigation/pool.png"
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
          icon: "/encrypt.png"
        },
        {
          name: "aside.systeminfo",
          path: "/Tools/System",
          icon: "/virus.png"
        },
        {
          name: "aside.data_handing",
          path: "/Tools/DataHanding",
          icon: "/data.png"
        },
        {
          name: "aside.memorandum",
          path: "/Tools/Memo",
          icon: "/memo.png"
        },
        {
          name: "aside.reverse",
          path: "/Tools/Reverse",
          icon: "/skull.png"
        },
        {
          name: "aside.associate_dictionary_generator",
          path: "/Tools/Thinkdict",
          icon: "/dict.png"
        },
        {
          name: "aside.aksk",
          path: "/Tools/AKSK",
          icon: "/wechat.png"
        },
      ]
    },
  ]

  export default MenuList