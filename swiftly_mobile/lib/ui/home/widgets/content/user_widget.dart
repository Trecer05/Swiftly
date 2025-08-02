import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/ui/home/widgets/call_made_widget.dart';
import 'package:swiftly_mobile/ui/home/widgets/content/avatar_widget.dart';

import '../../../../domain/user/models/user.dart';
import '../../../../providers/card_notifier_provider.dart';
import '../../../../providers/current_user_provider.dart';
import '../../../core/themes/colors.dart';
import '../../../core/themes/theme.dart';

class UserWidget extends ConsumerWidget {
  final User user;
  const UserWidget({super.key, required this.user});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final me = ref.watch(currentUserProvider);
    final cards =
        ref
            .watch(cardNotifierProvider)
            .cards
            .where((card) => card.userId == user.id)
            .toList();
    return Container(
      width: 180,
      height: 180,
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: me!.id == user.id ? AppColors.white64 : AppColors.white38,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              AvatarWidget(imageUrl: user.image),
              const SizedBox(width: 5),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '${user.name} ${user.lastName ?? ''}',
                      style: AppTextStyles.style13,
                      overflow: TextOverflow.ellipsis,
                      maxLines: 1,
                      softWrap: false,
                    ),
                    // if (user.role != null)
                    //   LabelItemWidget(labelItem: user.role!),
                  ],
                ),
              ),
            ],
          ),
          const Spacer(),
          if (me.id == user.id) const Align(alignment: Alignment.bottomRight, child: CallMadeWidget())
          else Row(
            crossAxisAlignment: CrossAxisAlignment.baseline,
            textBaseline: TextBaseline.alphabetic,
            children: [
              Text('${cards.length}', style: AppTextStyles.style6),
              const SizedBox(width: 5),
              const Text('Активных задач', style: AppTextStyles.style13),
            ],
          ),
        ],
      ),
    );
  }
}
